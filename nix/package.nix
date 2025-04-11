{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "todoist-cli";
  version = "unstable-2025-04-11";

  src = ../.;

  vendorHash = "sha256-Zig/atP/vbE8OPpbMzAa8dgbBiotOkh9bdkAMwBuy0U=";

  postFixup = ''
    mv $out/bin/todoist-cli $out/bin/todoist
  '';

  ldflags = [
    "-s"
    "-w"
  ];

  meta = {
    homepage = "https://github.com/CnTeng/todoist-cli";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [ CnTeng ];
    mainProgram = "todoist";
  };
}
