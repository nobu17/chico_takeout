import { auth, signIn } from "./Firebase";

export class AdminAuthService {
  private ids: string[];
  constructor() {
    this.ids = process.env.REACT_APP_ADMIN_IDS?.split(",") ?? [];
  }
  async signIn(email: string, password: string): Promise<AdminAuthResult> {
    try {
      const result = await signIn(auth, email, password);
      // check id are existed
      if (this.checkAdmin(email)) {
        return new AdminAuthResult(
          true,
          result.user.uid,
          result.user.email ?? "",
          ""
        );
      }
      return new AdminAuthResult(false, "", "", "管理ユーザではありません。");
    } catch (e: unknown) {
      let msg = "";
      if (e instanceof Error) {
        msg = e.message;
      }
      return new AdminAuthResult(false, "", "", msg);
    }
  }
  async signOut(): Promise<void> {
    await auth.signOut();
  }
  async onAuthStateChange(callback: (data: AdminAuthResult) => void) {
    auth.onAuthStateChanged((user) => {
      if (!user) {
        callback(new AdminAuthResult(false, "", "", ""));
      } else {
        // check id are existed
        if (this.checkAdmin(user.email)) {
          callback(new AdminAuthResult(true, user.uid, user.email ?? "", ""));
        } else {
          callback(
            new AdminAuthResult(false, "", "", "管理ユーザではありません。")
          );
        }
      }
    });
  }

  private checkAdmin(email: string | null): boolean {
    if (email === null) return false;
    return this.ids.includes(email);
  }
}

export class AdminAuthResult {
  constructor(
    public isSuccessful: boolean,
    public uid: string,
    public email: string,
    public errorMessage: string
  ) {}

  public hasError(): boolean {
    if (!this.isSuccessful && this.errorMessage !== "") {
      return true;
    }
    return false;
  }
}
