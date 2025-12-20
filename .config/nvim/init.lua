local set = vim.opt

vim.opt.packpath:append("~/.config/nvim/plugin")

vim.cmd('colorscheme shine')
vim.cmd('syntax on')
set.number = true
set.showcmd = true

set.hlsearch = true
set.incsearch = true
set.ignorecase = true
set.smartcase = true
set.wrapscan = true

set.virtualedit = 'onemore'
set.showmatch = true
set.autoindent = true
set.shiftwidth = 4
set.softtabstop = 4

set.clipboard:append({"unnamedplus"})

local _curfile = vim.fn.expand("%:r")
if _curfile == 'Makefile' then
    set.noexpandtab = true
    set.nosmarttab = true
end

set.writebackup = false
set.backup = false
set.swapfile = false

set.completeopt = 'menuone', 'preview'

set.scrolloff = 5

set.laststatus = 2
set.statusline = '%f%r%h%=%p'
vim.cmd('highlight StatusLine guifg=black guibg=darkcyan')

set.cursorline = true
vim.keymap.set('n', 'q:', '<Esc>', {silent = true, noremap = true})

-- autocmd! VimEnter * Ve | wincmd w
-- vim.cmd('highlight VertSplit ctermfg=gray ctermbg=gray')

set.showtabline = 2
vim.cmd('highlight TabLine guifg=white guibg=darkcyan')
vim.cmd('highlight TabLineSel guifg=black guibg=darkcyan')
vim.cmd('highlight TabLineFill guifg=white guibg=darkcyan')

local function transparent_background()
  local highlights = {
    "Normal", "NormalFloat", "NomalNC", "SignColumn",
    "MsgArea", "Pmenu", "TeleScopeBorder", "TelescopeNormal"
  }
  for _, hl in ipairs(highlights) do
    vim.api.nvim_set_hl(0, hl, { bg = "none" })
  end
end

transparent_background()

require('nvim-treesitter.configs').setup {
  ensure_installed = {},
  auto_install = false,
  highlight = {
    enable = true,
  },
  indent = {
    enable = true,
  }
}

-- require('comment_selector')
require('popup_menu_test')
require('buffer_tabline')
require('crucru')
require('stickey')
require('comment_bp')
vim.api.nvim_create_autocmd("FileType", {
    pattern = {"json", "scheme"},
    callback = function()
	vim.treesitter.stop()
    end,
})
