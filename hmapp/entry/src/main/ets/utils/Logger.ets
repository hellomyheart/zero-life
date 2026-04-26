export class Logger {
  private static domain: number = 0x0000;
  private static prefix: string = "[MyApp]";
  private static format: string = "%{public}s, %{public}s";
  private static isDebug: boolean = true;

  public static init(): void {
    this.isDebug = true;
  }

  public static debug(message: string, ...args: Object[]): void {
    if (this.isDebug) {
      hilog.debug(this.domain, this.prefix, this.format, message, args.length > 0 ? JSON.stringify(args) : "");
    }
  }

  public static info(message: string, ...args: Object[]): void {
    if (this.isDebug) {
      hilog.info(this.domain, this.prefix, this.format, message, args.length > 0 ? JSON.stringify(args) : "");
    }
  }

  public static warn(message: string, ...args: Object[]): void {
    if (this.isDebug) {
      hilog.warn(this.domain, this.prefix, this.format, message, args.length > 0 ? JSON.stringify(args) : "");
    }
  }

  public static error(message: string, ...args: Object[]): void {
    if (this.isDebug) {
      hilog.error(this.domain, this.prefix, this.format, message, args.length > 0 ? JSON.stringify(args) : "");
    }
  }

  public static fatal(message: string, ...args: Object[]): void {
    if (this.isDebug) {
      hilog.fatal(this.domain, this.prefix, this.format, message, args.length > 0 ? JSON.stringify(args) : "");
    }
  }

  public static isLoggable(): boolean {
    if (this.isDebug) {
      return true;
    }
    return false;
  }
}