self: nixosTest:
let
  stateVersion = "25.05";
  home-manager = builtins.fetchTarball {
    url = "https://github.com/nix-community/home-manager/archive/release-${stateVersion}.tar.gz";
    sha256 = "0d41gr0c89a4y4lllzdgmbm54h9kn9fjnmavwpgw0w9xwqwnzpax";
  };
in
nixosTest {
  name = "todoist-cli";
  nodes.machine = {
    imports = [ (import "${home-manager}/nixos") ];

    users.users.alice = {
      isNormalUser = true;
    };

    home-manager.sharedModules = [ self.homeModules.default ];
    home-manager.users.alice = {
      home.stateVersion = stateVersion;
      home.enableNixpkgsReleaseCheck = false;

      programs.todoist-cli.enable = true;
    };

    system.stateVersion = stateVersion;
  };

  testScript = ''
    machine.wait_for_unit("default.target")
    machine.succeed("su -- alice -c 'which todoist'")
  '';
}
