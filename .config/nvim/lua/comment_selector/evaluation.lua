local api = vim.api
local evaluation_file = "comments_evaluation.txt"

local M = {}
local evaluations = {}

local function load_evaluations()
  evaluations = {}
  local file = io.open(evaluation_file, "r")
  if not file then return end

  for line in file:lines() do
    local line_number, content, score = line:match("^(%d+)%s+(.+)%s+(-?%d+)$")
    if line_number and content and score then
      evaluations[tonumber(line_number)] = { content = content, score = tonumber(score) }
    end
  end
  file:close()
end

local function save_evaluations()
  local file = io.open(evaluation_file, "w")
  if not file then return end
  for line_number, data in pairs(evaluations) do
    file:write(string.format("%d %s %d\n", line_number, data.content, data.score))
  end
  file:close()
end

function M.update_evaluation(delta)
  local line_number = api.nvim_win_get_cursor(0)[1]
  local line_content = api.nvim_buf_get_lines(0, line_number - 1, line_number, false)[1]
  local cleaned = line_content:gsub("^%-%-%s*", "")

  evaluations[line_number] = evaluations[line_number] or { content = cleaned, score = 0 }
  evaluations[line_number].score = evaluations[line_number].score + delta
  print(string.format("評価更新: %s (評価: %d)", cleaned, evaluations[line_number].score))
  save_evaluations()
end

function M.show_comment_evaluations()
  load_evaluations()
  local quickfix = {}
  for line_number, data in pairs(evaluations) do
    table.insert(quickfix, {
      filename = api.nvim_buf_get_name(0),
      lnum = line_number,
      col = 1,
      text = string.format("評価: %d - %s", data.score, data.content),
    })
  end
  if #quickfix == 0 then
    print("評価されたコメントがありません")
    return
  end
  vim.fn.setqflist(quickfix, 'r')
  vim.cmd("copen")
end

function M.get_evaluation(line_number)
  return evaluations[line_number]
end

load_evaluations()
return M

