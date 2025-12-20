local function read_file_to_table(filepath)
  local file = io.open(filepath, "r")
  if not file then return {} end
  local lines = {}
  for line in file:lines() do table.insert(lines, line) end
  file:close()
  return lines
end

return {
  templates = read_file_to_table("./mat/templates.txt"),
  categories = read_file_to_table("./mat/categories.txt"),
  words_type = read_file_to_table("./mat/words_type.txt"),
  words_object = read_file_to_table("./matwords_object.txt"),
  words_directive = read_file_to_table("./mat/words_directive.txt"),
  words_situation = read_file_to_table("./mat/words_situation.txt"),
}

