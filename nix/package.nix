{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "todoist-cli";
  version = "unstable-2025-04-09";

  src = ../.;

  vendorHash = "sha256-xZ75hytf70hxCh2NpUxQFz8ZO86F0/lIgZ19rTlY5mI=";

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
