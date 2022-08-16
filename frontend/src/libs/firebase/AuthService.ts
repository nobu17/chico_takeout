import { User } from "firebase/auth";
import {
  auth,
  signIn,
  sendMail,
  sendResetMail,
  signUp,
  googleAuthProvider,
  signInRedirect,
  getRedirect,
} from "./Firebase";

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

  async signUp(email: string, password: string): Promise<SignUpResult> {
    try {
      const result = await signUp(auth, email, password);
      const mailSent = await this.checkEmailVerification(result.user);
      return new SignUpResult(
        true,
        result.user.uid,
        result.user.email ?? "",
        mailSent,
        ""
      );
    } catch (e: unknown) {
      let msg = "";
      if (e instanceof Error) {
        msg = e.message;
      }
      return new SignUpResult(false, "", "", false, msg);
    }
  }

  async sendPassResetMail(email: string): Promise<PassResetResult> {
    try {
      await sendResetMail(auth, email);
      return new PassResetResult(true, "");
    } catch (e: unknown) {
      let msg = "";
      if (e instanceof Error) {
        msg = e.message;
      }
      return new PassResetResult(false, msg);
    }
  }

  async signInWithGoogle() {
    const provider = new googleAuthProvider();
    await signInRedirect(auth, provider);
    // provider.addScope("email");
  }

  async getRedirectResult(callback: (data: AuthResult) => void) {
    const result = await getRedirect(auth);
    if (result !== null) {
      callback(new AuthResult(true, false, result.user.uid, result.user.email ?? "", ""))
      return;
    }
    callback(new AuthResult(false, false, "", "", ""));
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

export class SignUpResult {
  constructor(
    public isSuccessful: boolean,
    public uid: string,
    public email: string,
    public mailSent: boolean,
    public errorMessage: string
  ) {}

  public hasError(): boolean {
    if (!this.isSuccessful && this.errorMessage !== "") {
      return true;
    }
    return false;
  }
  public isUserAlreadyExists(): boolean {
    if (
      !this.isSuccessful &&
      this.errorMessage.includes("auth/email-already-in-use")
    ) {
      return true;
    }
    return false;
  }
}

export class PassResetResult {
  constructor(public isSuccessful: boolean, public errorMessage: string) {}

  public hasError(): boolean {
    if (!this.isSuccessful && this.errorMessage !== "") {
      return true;
    }
    return false;
  }
  public isUserNotExists(): boolean {
    if (
      !this.isSuccessful &&
      this.errorMessage.includes("auth/user-not-found")
    ) {
      return true;
    }
    return false;
  }
}
