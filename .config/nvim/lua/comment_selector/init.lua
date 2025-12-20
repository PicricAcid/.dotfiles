local eval = require("comment_selector.evaluation")
require("comment_selector.float_popup")

vim.api.nvim_create_user_command("SelectTemplate", require("comment_selector.comment_menu").comment_menu, {})
vim.api.nvim_create_user_command("CommentEvalList", eval.show_comment_evaluations, {})

vim.keymap.set('n', '<leader>b', function()
  eval.update_evaluation(1)
end, { noremap = true, silent = true })

vim.keymap.set('n', '<leader>q', function()
  eval.update_evaluation(-1)
end, { noremap = true, silent = true })

