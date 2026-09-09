{ pkgs, ... }: {
	home.packages = with pkgs; [
	arduino-cli
	cargo
	gopls
	lld
	llvm
	nkf
	nodejs
	qemu
	rustc
	sshfs
	terraform
	tree
	];
}
