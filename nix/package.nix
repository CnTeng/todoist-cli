{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "todoist-cli";
  version = "unstable-2025-04-30";

  src = ../.;

  vendorHash = "sha256-fM18bMW3Q5IoMhDovsbVkaJNR/McocC8ClfarFuOGpw=";

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
