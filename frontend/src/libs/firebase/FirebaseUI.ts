import { EmailAuthProvider } from "firebase/auth";
import { auth } from "./Firebase";
import * as firebaseui from "firebaseui";
import "firebaseui/dist/firebaseui.css";

export const uiConfig: firebaseui.auth.Config = {
  signInFlow: "popup",
  signInSuccessUrl: "/",
  signInOptions: [EmailAuthProvider.PROVIDER_ID],
  callbacks: {
    signInSuccessWithAuthResult: function (authResult) {
    //   if (authResult.additionalUserInfo.isNewUser) {
    //     authResult.user.sendEmailVerification();
    //   }
    //   console.log("authResult", authResult);
    //   //If this is a new user && his provider is password (username/password registeration) && his email is not verified,
    //   if (
    //     authResult.additionalUserInfo.isNewUser &&
    //     authResult.additionalUserInfo.providerId === "password" &&
    //     !authResult.user.emailVerified
    //   ) {
    //     //Send him the verification email, show him a toastr message, then apply a force logout
    //     //authResult.user.sendEmailVerification();
    //     //auth.signOut();
    //     alert(
    //       "アカウント作成が完了しました。メールアドレスに認証用のメールを送信しましたのでそちらから検証を行なってください。"
    //     );
    //   } else if (
    //     !authResult.additionalUserInfo.isNewUser &&
    //     authResult.additionalUserInfo.providerId === "password" &&
    //     !authResult.user.emailVerified
    //   ) {
    //     //Show him a taostr message that his email is not verified, then apply a force logout.
    //     const confirmResult = window.confirm(
    //       "Email認証が完了ていないようです。再度メールを送信しますか？"
    //     );
    //     if (confirmResult) {
    //       //authResult.user.sendEmailVerification();
    //     }
    //     //auth.signOut();
    //   } else {
    //     return true;
    //   }

      return false;
    },
  },
};
