{ pkgs, ... }: {
	home.packages = with pkgs; [
	nodejs
	sshfs
	tree
	];
}
