import * as React from "react";
import { NavLink as RouterLink } from "react-router-dom";
import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import Link from "@mui/material/Link";
import PageTitle from "../../components/parts/PageTitle";
import LoadingSpinner from "../../components/parts/LoadingSpinner";
import LoginForm, { LoginInput } from "../admin/parts/LoginForm";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../components/contexts/AuthContext";
import { useMessageDialog } from "../../hooks/UseMessageDialog";

export default function UserLogin() {
  const { showMessageDialog, renderDialog } = useMessageDialog();
  const navigate = useNavigate();
  const { signIn, loading, signInWithGoogle, signInWithTwitter } = useAuth();
  const handleSignIn = async (input: LoginInput) => {
    const result = await signIn(input.email, input.password);
    if (result.isSuccessful) {
      navigate("/my_page");
    } else if (!result.emailVerified) {
      await showMessageDialog("エラー", "メール認証が完了していません。先にメール認証を進めてください。");
    }
    else {
      await showMessageDialog("エラー", "ログインに失敗しました。");
    }
  };
  const handleSignInWithGoogle = async () => {
    await signInWithGoogle();
  };
  const handleSignInWithTwitter = async () => {
    await signInWithTwitter();
  };
  return (
    <>
      <PageTitle title="ログイン" />
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <LoginForm
            input={{ email: "", password: "" }}
            onSubmit={handleSignIn}
            onGoogleSubmit={handleSignInWithGoogle}
            onTwitterSubmit={handleSignInWithTwitter}
          ></LoginForm>
          <LoadingSpinner message="Loading..." isLoading={loading} />
        </Grid>
        <Grid item alignItems="center" xs={12}>
          <Typography align="center">
            <Link
              textAlign="center"
              underline="hover"
              color="error"
              component={RouterLink}
              to="/auth/sign_up"
            >
              新規登録はこちら
            </Link>
          </Typography>
        </Grid>
        <Grid item alignItems="center" xs={12}>
          <Typography align="center">
            <Link
              textAlign="center"
              underline="hover"
              color="primary"
              component={RouterLink}
              to="/auth/reset"
            >
              パスワードを忘れた場合はこちら
            </Link>
          </Typography>
        </Grid>
      </Grid>
      {renderDialog()}
    </>
  );
}
