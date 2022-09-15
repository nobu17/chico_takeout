import * as React from "react";
import Grid from "@mui/material/Grid";
import Alert from "@mui/material/Alert";
import LoadingSpinner from "../../components/parts/LoadingSpinner";
import PageTitle from "../../components/parts/PageTitle";
import { AuthService } from "../../libs/firebase/AuthService";
import SignUpForm, { SignUpInput } from "./parts/SignUpForm";
import { useNavigate } from "react-router-dom";
import { useMessageDialog } from "../../hooks/UseMessageDialog";

const service = new AuthService();

export default function UserSignUp() {
  const { showMessageDialog, renderDialog } = useMessageDialog();
  const navigate = useNavigate();
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<Error | undefined>(undefined);
  const [signUpFailedMessage, setSignUpFailedMessaged] =
    React.useState<string>("");
  const handleSignUp = async (input: SignUpInput) => {
    try {
      setSignUpFailedMessaged("");
      setError(undefined);
      setLoading(true);
      const result = await service.signUp(input.email, input.password);
      if (result.isSuccessful && result.mailSent) {
        await showMessageDialog(
          "",
          "確認用のメールを送信しました。メールのリンクから、本登録を行なってください。"
        );
        navigate("/");
        return;
      }
      if (result.isUserAlreadyExists()) {
        setSignUpFailedMessaged("既に登録されているメールアドレスです。");
      } else if (result.hasError()) {
        setSignUpFailedMessaged("エラーが発生しました。" + result.errorMessage);
      }
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };
  return (
    <>
      <PageTitle title="ユーザー登録" />
      <Grid container spacing={2}>
        <Grid item xs={12}>
          {error ? (
            <Alert severity="error">
              エラーが発生しました。お手数ですが時間をおいて再度お試しください。
            </Alert>
          ) : (
            <></>
          )}
          {signUpFailedMessage !== "" ? (
            <Alert severity="error">{signUpFailedMessage}</Alert>
          ) : (
            <></>
          )}
          <SignUpForm
            input={{ email: "", password: "" }}
            onSubmit={handleSignUp}
          ></SignUpForm>
        </Grid>
        <LoadingSpinner message="Loading..." isLoading={loading} />
      </Grid>
      {renderDialog()}
    </>
  );
}
