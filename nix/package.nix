{
  lib,
  buildGoModule,
  installShellFiles,
}:
buildGoModule {
  pname = "todoist-cli";
  version = "unstable-2025-05-23";

  src = ../.;

  vendorHash = "sha256-FMt3SoSkUlQoziuMRzSfJPgfc/jclkfbd5smdP6NIvU=";

  nativeBuildInputs = [ installShellFiles ];

  ldflags = [
    "-s"
    "-w"
  ];

  postInstall = ''
    mv $out/bin/todoist-cli $out/bin/todoist

    installShellCompletion --cmd todoist \
      --bash <($out/bin/todoist completion bash) \
      --zsh <($out/bin/todoist completion zsh) \
      --fish <($out/bin/todoist completion fish)
  '';

  meta = {
    description = "CLI client for Todoist";
    homepage = "https://github.com/CnTeng/todoist-cli";
    license = lib.licenses.mit;
    mainProgram = "todoist";
    maintainers = with lib.maintainers; [ CnTeng ];
    platforms = lib.platforms.all;
  };
}
