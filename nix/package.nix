{
  lib,
  buildGoModule,
  installShellFiles,
}:
buildGoModule {
  pname = "todoist-cli";
  version = "unstable-2025-08-27";

  src = ../.;

  vendorHash = "sha256-R4H1auwlOkbfOCq6BuRJJYWcl/sQAPaQ3orOkLWjkAA=";

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
