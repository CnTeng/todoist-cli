{
  description = "A Todoist CLI client";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";

    git-hooks-nix = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    treefmt = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ self, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      imports = [
        inputs.git-hooks-nix.flakeModule
        inputs.treefmt.flakeModule
      ];

      flake.overlays.default = final: prev: {
        todoist-cli = final.callPackage ./nix/package.nix { };
      };

      flake.homeManagerModules.default = import ./nix/hm-module.nix self;

      perSystem =
        {
          config,
          pkgs,
          system,
          ...
        }:
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = [ self.overlays.default ];
          };

          packages = {
            default = config.packages.todoist-cli;
            todoist-cli = pkgs.todoist-cli;
          };

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              go
              config.packages.todoist-cli
              config.treefmt.build.wrapper
            ];
            CGO_ENABLED = "0";
            shellHook = config.pre-commit.installationScript;
          };

          pre-commit.settings.hooks = {
            commitizen.enable = true;
            treefmt.enable = true;
          };

          treefmt.programs = {
            gofumpt.enable = true;
            nixfmt.enable = true;
            prettier.enable = true;
          };
        };
    };
}
