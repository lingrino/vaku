use std::process::ExitCode;

mod cli;

fn main() -> ExitCode {
    let version = env!("CARGO_PKG_VERSION");
    let args: Vec<String> = std::env::args().skip(1).collect();
    let code = cli::execute(
        version,
        &args,
        &mut std::io::stdout(),
        &mut std::io::stderr(),
    );
    ExitCode::from(code)
}
