local client = vim.lsp.start_client {
    name = "go-lsp",
    cmd = { "/Users/user/code/go-lsp/bin/golsp" },
    on_attach = on_attach,
}

if not client then 
    vim.notify "client initialization failed"
    return
end

vim.api.nvim_create_autocmd("FileType", {
    pattern = { "go", "markdown" },
    callback = function()
        vim.lsp.buf_attach_client(0, client)
    end,
})