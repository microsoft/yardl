#include <filesystem>

#include "generated/binary/protocols.h"

void walkTree(sketch::BinaryTree const& value, std::function<void(sketch::BinaryTree const&)> const& callback) {
  if (value.left) {
    walkTree(*value.left, callback);
  }

  callback(value);

  if (value.right) {
    walkTree(*value.right, callback);
  }
}

void walkTree(std::unique_ptr<sketch::BinaryTree> const& value, std::function<void(sketch::BinaryTree const&)> const& callback) {
  if (!value) {
    return;
  }

  walkTree(value->left, callback);
  callback(*value);
  walkTree(value->right, callback);
}

void insertTree(std::unique_ptr<sketch::BinaryTree>& root, int value) {
  if (!root) {
    root = std::make_unique<sketch::BinaryTree>();
    root->value = value;
    return;
  }

  if (value < root->value) {
    if (root->left) {
      insertTree(root->left, value);
    } else {
      root->left = std::make_unique<sketch::BinaryTree>();
      root->left->value = value;
    }
  } else {
    if (root->right) {
      insertTree(root->right, value);
    } else {
      root->right = std::make_unique<sketch::BinaryTree>();
      root->right->value = value;
    }
  }
}

void readDirectory(sketch::Directory& dir, std::filesystem::path const& path) {
  for (auto const& dir_entry : std::filesystem::directory_iterator(path)) {
    if (dir_entry.is_directory()) {
      auto sub_dir = std::make_unique<sketch::Directory>();
      sub_dir->name = dir_entry.path().filename().string();
      readDirectory(*sub_dir, dir_entry.path());
      dir.entries.emplace_back(std::move(sub_dir));
    } else {
      sketch::File file;
      file.name = dir_entry.path().filename().string();
      // file->data = std::vector<uint8_t>(std::filesystem::file_size(dir_entry.path()));
      // std::ifstream file_stream(dir_entry.path(), std::ios::binary);
      // file_stream.read(reinterpret_cast<char*>(file->data.data()), file->data.size());
      dir.entries.emplace_back(std::move(file));
    }
  }
}

void writeDirectoryEntry(sketch::DirectoryEntry const& entry, std::filesystem::path const& path) {
  (void)(path);

  if (std::holds_alternative<sketch::File>(entry)) {
    auto const& file = std::get<sketch::File>(entry);
    auto file_path = path / file.name;
    std::cerr << "touch " << file_path << std::endl;
    // std::ofstream file_stream(file_path, std::ios::binary);
    // file_stream.write(reinterpret_cast<char const*>(file.data.data()), file.data.size());
  } else {
    auto const& sub_dir = std::get<std::unique_ptr<sketch::Directory>>(entry);
    // std::filesystem::create_directory(sub_dir->name);
    auto dir_path = path / sub_dir->name;
    std::cerr << "mkdir " << dir_path << std::endl;
    for (auto const& sub_entry : sub_dir->entries) {
      writeDirectoryEntry(sub_entry, dir_path);
    }
  }
}

int main(void) {
  std::unique_ptr<sketch::BinaryTree> root;
  for (int i = 0; i < 32; i++) {
    insertTree(root, (rand() % 100) - 50);
  }
  walkTree(*root, [](sketch::BinaryTree const& node) {
    std::cerr << node.value << " ";
  });
  std::cerr << std::endl;

  std::stringstream output;
  sketch::binary::MyProtocolWriter writer(output);
  writer.WriteTree(*root);
  writer.WritePtree(root);

  sketch::LinkedList<std::string> list;
  list.value = "Hello";
  list.next = std::make_unique<sketch::LinkedList<std::string>>();
  list.next->value = "World";
  list.next->next = std::make_unique<sketch::LinkedList<std::string>>();
  list.next->next->value = "!!!";
  writer.WriteList(std::move(list));

  sketch::Directory cwd;
  readDirectory(cwd, std::filesystem::current_path());
  for (auto const& entry : cwd.entries) {
    writer.WriteCwd(entry);
  }
  writer.EndCwd();

  writer.Close();

  /////////////////////////////////////////////////////////////////////////////

  std::stringstream input(output.str());
  sketch::binary::MyProtocolReader reader(input);

  sketch::BinaryTree result;
  reader.ReadTree(result);
  walkTree(result, [](sketch::BinaryTree const& node) {
    std::cerr << node.value << " ";
  });
  std::cerr << std::endl;

  std::unique_ptr<sketch::BinaryTree> presult;
  reader.ReadPtree(presult);
  walkTree(presult, [](sketch::BinaryTree const& node) {
    std::cerr << node.value << " ";
  });
  std::cerr << std::endl;

  std::optional<sketch::LinkedList<std::string>> olist;
  reader.ReadList(olist);
  if (olist) {
    auto list = std::make_unique<sketch::LinkedList<std::string>>();
    *list = std::move(*olist);
    while (list) {
      std::cerr << list->value << " ";
      list = std::move(list->next);
    }
    std::cerr << std::endl;
  }

  sketch::DirectoryEntry entry;
  while (reader.ReadCwd(entry)) {
    writeDirectoryEntry(entry, ".");
  }

  reader.Close();

  return 0;
}
