#include "src/argparser.h"
#include "src/deque.h"
#include "src/diff.h"
#include "src/gitobjects.h"
#include "src/hashmap.h"
#include "src/helpers.h"
#include "src/index.h"
#include "src/objectstore.h"
#include "src/refs.h"
#include "src/vector.h"
#include <exception>
#include <filesystem>
#include <fstream>
#include <functional>
#include <iostream>
#include <ostream>
#include <stdexcept>
#include <string>
#include <curl/curl.h>

int main(int argc, char *argv[])
{
  // Argument parser for command line interface.
  ArgParser parser(argv[0], "Jit Version Control System.");
  parser.add_command("help", "Show this help message").set_callback([&]()
                                                                    { std::cout << parser.help_message() << std::endl; });

  // Initializing a jit repository.
  parser.add_command("init", "Initialize a repository").set_callback([&]()
                                                                     {
    std::filesystem::path repo = "./.jit";
    if (std::filesystem::exists(repo))
      throw std::runtime_error(".jit directory already exists");
    Refs refs(repo / "refs", repo / "HEAD");
    refs.updateHead("main"); });

  // Adding files or directories to the staging area.
  std::string addPath;
  parser.add_command("add", "Add file to the staging area")
      .set_callback([&]()
                    {
        std::filesystem::path repo = repoRoot();
        ObjectStore store(repo / "objects");
        IndexStore index(repo / "index", store);
        index.add(addPath);
        index.save(); })
      .add_argument(addPath, "File Path", "");

  // Commiting changes into the repository.
  std::string commitMessage;
  parser.add_command("commit", "Add file to the staging area")
      .set_callback([&]()
                    {
        std::filesystem::path repo = repoRoot();
        ObjectStore store(repo / "objects");
        IndexStore index(repo / "index", store);
        Refs refs(repo / "refs", repo / "HEAD");

        std::string current = refs.resolve("HEAD");

        Tree commitTree = index.writeTree();
        store.store(&commitTree);

        Commit *newCommit =
            new Commit(commitMessage, "pharoak", commitTree.getHash());

        if (current != "")
          newCommit->addParentHash(current);

        if (refs.getMergeHead() != "") {
          std::string h = refs.resolve(refs.getMergeHead());
          newCommit->addParentHash(h);
          refs.updateMergeHead("");
        }

        store.store(newCommit);
        if (refs.isHeadBranch())
          refs.updateRef(refs.getHead(), newCommit->getHash()); })
      .add_option(commitMessage, "-m,--message",
                  "Must be between double quotations.");

  // Shows the logs of the commits stored in the repository.
  parser.add_command("log", "Display the log of the commits")
      .set_callback([&]()
                    {
        std::filesystem::path repo = repoRoot();
        ObjectStore store(repo / "objects");
        Refs refs(repo / "refs", repo / "HEAD");

        std::string lastCommit = refs.resolve("HEAD");
        if (lastCommit == "") {
          std::cout << "your current branch does not have any commits yet."
                    << std::endl;
          return;
        }

        Vector<Pair<std::string, std::string>> log;
        store.retrieveLog(lastCommit, log);
        for (auto &l : log)
          std::cout << l.second << std::endl; });

  // Computes the difference between different commits and files.
  std::string filePath1, filePath2;
  parser.add_command("diff", "Computes the differences between files")
      .set_callback([&]()
                    {
        std::filesystem::path repo = repoRoot();
        ObjectStore store(repo / "objects");
        IndexStore index(repo / "index", store);
        Refs refs(repo / "refs", repo / "HEAD");
        std::string current = refs.resolve("HEAD");

        if (filePath1.empty() && filePath2.empty()) {
          // Diff between current commit and Staging area
          Tree commitTree = index.writeTree();
          store.store(&commitTree);

          HashMap<std::string, std::string> Blobs;

          Vector<TreeEntry> currentTree = commitTree.getEntries();
          for (int i = 0; i < currentTree.size(); i++) {
            if (currentTree[i].type == "blob") {
              Blobs.set(currentTree[i].name, currentTree[i].hash);
            } else if (currentTree[i].type == "tree") {
              GitObject *obj = store.retrieve(currentTree[i].hash);
              if (Tree *tree = dynamic_cast<Tree *>(obj)) {
                Vector<TreeEntry> subentries = tree->getEntries();
                for (int j = 0; j < subentries.size(); j++) {
                  currentTree.push_back(subentries[j]);
                }
              }
            }
          }

          Vector<std::string> results;
          GitObject *commitObj = store.retrieve(current);
          if (Commit *commit = dynamic_cast<Commit *>(commitObj)) {
            std::string tree = commit->getTreeHash();
            GitObject *treeObj = store.retrieve(tree);
            if (Tree *tree = dynamic_cast<Tree *>(treeObj)) {
              Vector<TreeEntry> retrievedTree = tree->getEntries();
              for (int i = 0; i < retrievedTree.size(); i++) {
                if (retrievedTree[i].type == "blob") {
                  std::string blobHash = Blobs.get(retrievedTree[i].name);
                  results.push_back("---" + retrievedTree[i].name + "---");
                  if (retrievedTree[i].hash == blobHash || blobHash == "") {
                    results.push_back("No Difference Found");
                  } else {
                    Vector<std::string> file1;
                    Vector<std::string> file2;
                    GitObject *blobObj1 = store.retrieve(retrievedTree[i].hash);
                    GitObject *blobObj2 = store.retrieve(blobHash);
                    if (Blob *blob1 = dynamic_cast<Blob *>(blobObj1)) {
                      if (Blob *blob2 = dynamic_cast<Blob *>(blobObj2)) {
                        file1 = split(blob1->getContent(), '\n');
                        file2 = split(blob2->getContent(), '\n');
                      }
                    }

                    Vector<std::string> differences = diff(file1, file2);
                    for (int k = 0; k < differences.size(); k++) {
                      results.push_back(differences[k]);
                    }
                  }
                } else if (retrievedTree[i].type == "tree") {
                  GitObject *obj = store.retrieve(retrievedTree[i].hash);
                  if (Tree *tree = dynamic_cast<Tree *>(obj)) {
                    Vector<TreeEntry> subentries = tree->getEntries();
                    for (int j = 0; j < subentries.size(); j++) {
                      retrievedTree.push_back(subentries[j]);
                    }
                  }
                }
              }
            }
          }

          std::cout << "File Differences:" << "\n"
                    << "================" << "\n";
          for (auto str : results) {
            std::cout << str << "\n";
          }
        } else if (!filePath2.empty()) {
          // Diff between a file and another
          std::cout << "File Differences:" << "\n"
                    << "================" << "\n";

          std::ifstream file1(filePath1), file2(filePath2);
          Vector<std::string> lines1, lines2;
          std::string line;

          if (!file1.is_open()) {
            throw std::runtime_error("Cannot open file: " + filePath1);
          }
          while (std::getline(file1, line)) {
            lines1.push_back(line);
          }
          if (!file2.is_open()) {
            throw std::runtime_error("Cannot open file: " + filePath2);
          }
          while (std::getline(file2, line)) {
            lines2.push_back(line);
          }
          Vector<std::string> result = diff(lines1, lines2);
          for (const auto &line : result)
            std::cout << line << "\n";
        } else {
          // Diff between a commit and current commit
          HashMap<std::string, std::string> Blobs;

          GitObject *headObj = store.retrieve(current);
          if (Commit *commit = dynamic_cast<Commit *>(headObj)) {
            std::string tree = commit->getTreeHash();
            GitObject *treeObj = store.retrieve(tree);
            if (Tree *tree = dynamic_cast<Tree *>(treeObj)) {
              Vector<TreeEntry> currentTree = tree->getEntries();
              for (int i = 0; i < currentTree.size(); i++) {
                if (currentTree[i].type == "blob") {
                  Blobs.set(currentTree[i].name, currentTree[i].hash);
                } else if (currentTree[i].type == "tree") {
                  GitObject *obj = store.retrieve(currentTree[i].hash);
                  if (Tree *tree = dynamic_cast<Tree *>(obj)) {
                    Vector<TreeEntry> subentries = tree->getEntries();
                    for (int j = 0; j < subentries.size(); j++) {
                      currentTree.push_back(subentries[j]);
                    }
                  }
                }
              }
            }
          }

          Vector<std::string> results;

          GitObject *commitObj = store.retrieve(filePath1);
          if (commitObj == nullptr) {
            return void(std::cout << "Hash does not exits\n");
          }
          if (Commit *commit = dynamic_cast<Commit *>(commitObj)) {
            std::string tree = commit->getTreeHash();
            GitObject *treeObj = store.retrieve(tree);
            if (Tree *tree = dynamic_cast<Tree *>(treeObj)) {
              Vector<TreeEntry> retrievedTree = tree->getEntries();
              for (int i = 0; i < retrievedTree.size(); i++) {
                if (retrievedTree[i].type == "blob") {
                  std::string blobHash = Blobs.get(retrievedTree[i].name);
                  results.push_back("---" + retrievedTree[i].name + "---");
                  if (retrievedTree[i].hash == blobHash || blobHash == "") {
                    results.push_back("No Difference Found");
                  } else {
                    Vector<std::string> file1;
                    Vector<std::string> file2;

                    GitObject *blobObj1 = store.retrieve(retrievedTree[i].hash);
                    GitObject *blobObj2 = store.retrieve(blobHash);
                    if (Blob *blob1 = dynamic_cast<Blob *>(blobObj1)) {
                      if (Blob *blob2 = dynamic_cast<Blob *>(blobObj2)) {
                        file1 = split(blob1->getContent(), '\n');
                        file2 = split(blob2->getContent(), '\n');
                      }
                    }

                    Vector<std::string> differences = diff(file1, file2);
                    for (int k = 0; k < differences.size(); k++) {
                      results.push_back(differences[k]);
                    }
                  }
                } else if (retrievedTree[i].type == "tree") {
                  GitObject *obj = store.retrieve(retrievedTree[i].hash);
                  if (Tree *tree = dynamic_cast<Tree *>(obj)) {
                    Vector<TreeEntry> subentries = tree->getEntries();
                    for (int j = 0; j < subentries.size(); j++) {
                      retrievedTree.push_back(subentries[j]);
                    }
                  }
                }
              }
            }
          }

          std::cout << "File Differences:" << "\n"
                    << "================" << "\n";
          for (auto str : results) {
            std::cout << str << "\n";
          }
        } })
      .add_argument(filePath1, "", "", false)
      .add_argument(filePath2, "", "", false);

  // Checks status of the the staging area
  parser
      .add_command(
          "status",
          "Shows the tracked and untracked files in the working repository.")
      .set_callback([&]()
                    {
        std::filesystem::path repo = repoRoot();
        std::filesystem::path wd = repo.parent_path();
        ObjectStore store(repo / "objects");
        IndexStore index(repo / "index", store);

        enum class Status { Clean, NewFile, Modified, Deleted };

        // TODO: improve array search
        HashMap<std::string, Status> status;

        Vector<std::string> untracked;
        for (auto it = std::filesystem::recursive_directory_iterator(wd);
             it != std::filesystem::end(it); it++) {
          const auto &entry = *it;
          if (entry.path().filename().string().rfind(".jit", 0) == 0) {
            it.disable_recursion_pending();
            continue;
          }

          if (std::filesystem::is_regular_file(entry))
            untracked.push_back(pathString(entry));
        }

        for (auto &[path, hash] : index.getEntries()) {
          bool found = false;
          for (auto &p : untracked) {
            if (p == path) {
              found = true;
              break;
            }
          }
          if (found) {
            Blob untrackedBlob = Blob(readFile(wd / path));
            std::string untrackedHash = untrackedBlob.getHash();
            if (untrackedHash != hash)
              status[path] = Status::Modified;
          } else {
            status[path] = Status::Deleted;
          }
        }

        for (auto &path : untracked) {
          bool found = false;
          for (auto &[p, _] : index.getEntries()) {
            if (p == path) {
              found = true;
              break;
            }
          }
          if (!found)
            status[path] = Status::NewFile;
        }

        bool clean = true;
        // TODO: sort by status
        for (auto [path, status] : status) {
          switch (status) {
          case Status::NewFile:
            std::cout << "new file: " << path << std::endl;
            clean = false;
            break;
          case Status::Modified:
            std::cout << "modified: " << path << std::endl;
            clean = false;
            break;
          case Status::Deleted:
            std::cout << "deleted: " << path << std::endl;
            clean = false;
            break;
          }
        }
        if (clean) {
          std::cout << "Working tree clean." << std::endl;
        } });

  std::string target;
  parser.add_command("checkout", "Switches to a branch or to a commit")
      .set_callback([&]()
                    {
        std::filesystem::path repo = repoRoot();
        std::filesystem::path wd = repo.parent_path();
        ObjectStore store(repo / "objects");
        IndexStore index(repo / "index", store);
        Refs refs(repo / "refs", repo / "HEAD");

        std::string hash =
            refs.isBranch(target) ? refs.resolve(target) : target;
        if (Commit *c = dynamic_cast<Commit *>(store.retrieve(hash))) {
          store.reconstruct(c->getTreeHash(), wd);
          index.readTree(wd, c->getTreeHash());
          refs.updateHead(target);
          index.save();
        } else {
          std::cout << "no such branch or commit '" << target << "'"
                    << std::endl;
        } })
      .add_argument(target, "Commit hash", "");

  std::string branchName;
  parser.add_command("branch", "Creates a branch in the working directory.")
      .set_callback([&]()
                    {
        std::filesystem::path repo = repoRoot();
        Refs refs(repo / "refs", repo / "HEAD");
        if (branchName.empty()) {
          Vector<std::string> branches = refs.getRefs();
          for (auto &b : branches)
            std::cout << (b == refs.getHead() ? "+" : " ") << b << std::endl;
          return;
        }
        refs.updateRef(branchName, "HEAD"); })
      .add_argument(branchName, "", "", false);

  parser.add_command("merge", "Merge two branches together.")
      .set_callback([&]()
                    {
        std::filesystem::path repo = repoRoot();
        std::filesystem::path wd = repo.parent_path();
        ObjectStore store(repo / "objects");
        IndexStore index(repo / "index", store);
        Refs refs(repo / "refs", repo / "HEAD");

        auto resolveCommit = [&](std::string s) {
          std::string commitHash = refs.resolve(s);

          Commit *commit = nullptr;
          if (Commit *c = dynamic_cast<Commit *>(store.retrieve(commitHash)))
            commit = c;
          return commit;
        };

        Commit *ourHead = resolveCommit(refs.getHead());
        Commit *otherHead = resolveCommit(branchName);

        if (ourHead == nullptr) {
          std::cout << "head is detached" << std::endl;
          return;
        }

        if (otherHead == nullptr) {
          std::cout << "unknown branch '" << branchName << "'" << std::endl;
          return;
        }

        // fast-forward merge
        Deque<Commit *> dq;
        dq.push_back(otherHead);
        while (!dq.empty()) {
          Commit *cur = dq.front();
          dq.pop_front();
          if (cur->getHash() == ourHead->getHash()) {
            // otherHead->getHash() is different for some reason
            refs.updateRef(refs.getHead(), refs.resolve(branchName));
            store.reconstruct(otherHead->getTreeHash(), wd);
            std::cout << "performed fast-forward merge" << std::endl;
            return;
          }
          for (auto &h : cur->getParentHashes()) {
            Commit *next = resolveCommit(h);

            if (next != nullptr)
              dq.push_back(next);
          }
        }

        // divergent branches
        using Path = std::filesystem::path;
        using H = HashMap<Path, Blob *>;
        HashMap<Path, Blob *> ourBlobs, otherBlobs;

        std::function<void(std::string, Path, HashMap<Path, Blob *> &)>
            collect = [&](std::string hash, std::filesystem::path p,
                          HashMap<Path, Blob *> &h) {
              GitObject *obj = store.retrieve(hash);
              if (Blob *b = dynamic_cast<Blob *>(obj)) {
                h[p] = b;
              } else if (Tree *t = dynamic_cast<Tree *>(obj)) {
                for (auto entry : t->getEntries())
                  collect(entry.getHash(), p / entry.name, h);
              }
            };

        collect(ourHead->getTreeHash(), wd, ourBlobs);
        collect(otherHead->getTreeHash(), wd, otherBlobs);

        for (auto [path, blobp] : otherBlobs) {
          Blob *our = ourBlobs[path];
          if (our == nullptr) {
            // new incoming file
            store.reconstruct(blobp->getHash(), path);
          } else {
            Vector<std::string> d = diff(split(our->getContent(), '\n'),
                                         split(blobp->getContent(), '\n'));
            std::string newContent;

            int marker = -1;
            auto advanceMarker = [&]() {
              if (marker == 0)
                newContent += "<<<<<<<<< HEAD\n";
              if (marker == 1)
                newContent += "========\n";
              if (marker == 2)
                newContent += ">>>>>>>>> " + branchName + '\n';
              marker++;
              marker %= 3;
            };
            for (auto l : d) {
              int x = std::string(" -+").find(l[0]);
              while (marker != x)
                advanceMarker();
              newContent += l.substr(1) + '\n';
            }
            while (marker != 0)
              advanceMarker();

            std::ofstream file(path);
            file << newContent;
          }

          refs.updateMergeHead(refs.resolve(branchName));
        } })
      .add_argument(branchName, "branch name", "branch to merge");

  // Remote Repository
  /*
  Commands:
  User Authentication
  -- auth login
  -- auth register
  -- profile // get the current user profile
  -- profile -u,--username <username> // get profile of the specified user
  Repository Management
  -- repo // print current
  -- repo -u,--username <username> -r <reponame>// Show all repos for the current user
  -- repo --all // show all repos for the current user 
  -- repo create // create a repo for the current logged in user 
  -- repo delete // delete a repo for the cuurrent logged in user
  -- repo -l <link> // print the current repo information
  Repository Authorization
  -- grant <user> <link> // for current repo
  -- revoke <user> <link> // for current repo
  -- remote // print current remote origin
  -- remote add <link> // One Repo One Origin
  -- remote remove // remove current origin
  File Upload & Download
  -- push // upload the commit history to the remote repo 
  -- pull // download the commit history from the remote repo
  */

  std::string username, password, email_address, full_name,repourl,reponame;
  bool all;
  // Authorization
  parser.add_command("jithub login", "login to jithub remote server").set_callback([&]() {

  })
  .add_option(username, "-u,--username", "username of the user in the server", true).add_option(password, "-p,--password", "password for the user in the server", true);

  parser.add_command("jithub register", "register to jithub remote server").set_callback([&]() {

  })
  .add_option(username, "-u,--username", "username for the user in the server", true)
  .add_option(password, "-p,--password", "password for the user in the server", true)
  .add_option(full_name, "-n,--name", "full name for the user in the server", true)
  .add_option(email_address, "-e,--email", "email address for the user in the server", true);

  parser.add_command("profile","show profile for the specified user").set_callback([&]() {
      if (username == "") {

      } else  {

      }

  }).add_option(username,"-u,--username","username");
  
  
  //Repository Management
  parser.add_command("repo","repository management according to the flags specified").set_callback([&]() {

  })
  .add_option(username,"-u,--username","show repos for a user")
  .add_option(reponame,"-r,--repo","show repo details for a specified user")
  .add_option(all,"--all","show all repos")
  .add_option(repourl,"-r,--url","show repo data for the specified url");

  parser.add_command("repo create","create a repository for the current user").set_callback([&]() {

  });

  parser.add_command("repo delete","delete a repository for the current user").set_callback([&]() {

  }).add_argument(reponame,"reponame","repo needs to be deleted for the current user");

  parser.add_command("grant","grant access on the repository").set_callback([&](){

  })
  .add_argument(username,"username","target username for granting")
  .add_argument(reponame,"reponame","target reponame for granting, must be the current logged in user and the owner");

  parser.add_command("revoke","revoke access on the repository").set_callback([&]() {

  })
  .add_argument(username,"username","target username for revoking")
  .add_argument(reponame,"reponame","target reponame for revoking, must be the current logged in user and the owner");

  // Remote Origin Management  
  parser.add_command("remote","Show remote origin for this repository").set_callback([&]() {
    
  });
  
  parser.add_command("remote add", "add remote server to the current repository").set_callback([&]() {
    
  })
  .add_argument(repourl, "repo_url", "repository url in the form \"jithub.com/repoowner/reponame\"");
  
  parser.add_command("remote remove", "remove remote server from the current repository").set_callback([&]() {
    
  });
    
  // Push & Pull Service
  parser.add_command("push","push the current working tree to the remote repository").set_callback([&]() {

  });

  parser.add_command("pull", "pull the working tree from the remote repository").set_callback([&]() {

  });
  
  try
  {
    parser.parse(argc, argv);
  }
  catch (const std::exception &e)
  {
    std::cout << e.what() << std::endl;
  }

  return 0;
}