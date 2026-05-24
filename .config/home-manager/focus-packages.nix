{ pkgs, ... }: {
	home.packages = with pkgs; [
	arduino-cli
	gopls
	lld
	llvm
	nkf
	nodejs
	qemu
	sshfs
	tree
	];
}
