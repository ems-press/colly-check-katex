{
  description = "Basic go environment";
  inputs.nixpkgs.url = github:NixOS/nixpkgs;

  outputs = { self, nixpkgs }:
  let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
    lib = pkgs.lib;
  in {
    devShell.${system} = pkgs.mkShell rec {
      buildInputs = with pkgs; [
        go
        gopls
        gotools
      ];
      # LD_LIBRARY_PATH = "${lib.makeLibraryPath buildInputs}";
    };
  };
}
