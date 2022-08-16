import { User } from "firebase/auth";
import { auth, signIn, sendMail } from "./Firebase";

export class AuthService {
  private adminIds: string[];
  constructor() {
    this.adminIds = process.env.REACT_APP_ADMIN_IDS?.split(",") ?? [];
  }
  async signIn(email: string, password: string): Promise<AuthResult> {
    try {
      const result = await signIn(auth, email, password);
      // check id are existed
      const isAdmin = await this.checkAdmin(result.user);
      return new AuthResult(
        true,
        isAdmin,
        result.user.uid,
        result.user.email ?? "",
        ""
      );
    } catch (e: unknown) {
      let msg = "";
      if (e instanceof Error) {
        msg = e.message;
      }
      return new AuthResult(false, false, "", "", msg);
    }
  }
  async signOut(): Promise<void> {
    await auth.signOut();
  }
  async onAuthStateChange(callback: (data: AuthResult) => void) {
    auth.onAuthStateChanged(async (user) => {
      if (!user) {
        callback(new AuthResult(false, false, "", "", ""));
      } else {
        const isAdmin = await this.checkAdmin(user);
        if (!isAdmin) {
          const checkRes = await this.checkEmailVerification(user);
          if (!checkRes) {
            callback(
              new AuthResult(false, isAdmin, user.uid, user.email ?? "", "")
            );
            await this.signOut();
            return;
          }
        }
        callback(new AuthResult(true, isAdmin, user.uid, user.email ?? "", ""));
      }
    });
  }

  private mailSent: boolean = false;

  private async checkEmailVerification(user: User): Promise<boolean> {
    if (user.emailVerified) {
      return true;
    }
    if (!this.mailSent) {
      this.mailSent = true;
      alert(
        "検証用のメールを送信しました。メールからユーザー登録を行なってください。"
      );
      await sendMail(user);
    }

    return false;
  }

  private async checkAdmin(user: User | null): Promise<boolean> {
    if (user === null) return false;
    const token = await user.getIdTokenResult();
    // console.log("token", token);
    const role = token.claims["role"];
    return role === "Admin";
  }
}

export class AuthResult {
  constructor(
    public isSuccessful: boolean,
    public isAdmin: boolean,
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
