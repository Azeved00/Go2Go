{
    description = "P2P in Go";

    inputs = {
        nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
	};

    outputs = { self, ... } @ inputs: 
    let
        pkgs = import inputs.nixpkgs { inherit system; };
        ROOT = let p = builtins.getEnv "PWD"; in if p == "" then self else p;
        name = "P2P In Go";
        system = "x86_64-linux";
    in {
        devShells."${system}".default = pkgs.mkShell {
            inherit name ROOT;

            buildInputs = with pkgs; [ go gopls ];

            shellHook = ''
                alias trab1=$ROOT/trab1/runner.sh
                alias trab2=$ROOT/trab1/runner.sh
            '';
        };
    };
}
