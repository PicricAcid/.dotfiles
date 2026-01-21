local M = {}

M.searched_dir = ""
M.searched_filetype = ""
M.find_dir_tmp = nil

function M.find_dir(filetype, search_dir)
    if
	not M.find_dir_tmp
	or not string.find(search_dir, M.searched_dir, 1, true)
	or not string.find(filetype, M.searched_filetype, 1, true)
    then
	M.searched_dir = search_dir
	M.searched_filetype = filetype

	local find_sh = string.format(
	    'find %s -name %s -type f ! -path "*.git/*"',
	    search_dir,
	    filetype
	)

	local result = vim.fn.system(find_sh)
	M.find_dir_tmp = vim.fn.tempname()
	vim.fn.writefile(vim.split(result, "\n"), M.find_dir_tmp)
    end
end

function M.grep(text)
    if not M.find_dir_tmp then
	return
    end

    local grep_sh = string.format(
	'cat %s | xargs grep -n %s /dev/null',
	M.find_dir_tmp,
	vim.fn.shellescape(text)
    )

    local output = vim.fn.system(grep_sh)
    vim.fn.setqflist({}, 'r', { lines = vim.split(output, "\n") })
    vim.cmd("copen")
end

function M.grep_all(opts)
    local args = opts.fargs
    local argc = #args

    if argc < 1 then
	M.find_dir('"*.*"', "./")
	M.grep(vim.fn.expand("<cword>"))
    elseif argc == 1 then
	M.find_dir('"*.*"', "./")
	M.grep(args[1])
    elseif argc ==2 then
	M.find_dir(args[2], args[3])
	M.grep(args[1])
    end
end

return M
