{ pkgs, ... }: {
	home.packages = with pkgs; [
	gopls
	nodejs
	sshfs
	tree
	];
}
