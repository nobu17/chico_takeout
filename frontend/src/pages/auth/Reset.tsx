import * as React from "react";
import Grid from "@mui/material/Grid";
import Alert from "@mui/material/Alert";
import LoadingSpinner from "../../components/parts/LoadingSpinner";
import PageTitle from "../../components/parts/PageTitle";
import { AuthService } from "../../libs/firebase/AuthService";
import { useNavigate } from "react-router-dom";
import ResetForm, { ResetInput } from "./parts/ResetForm";

const service = new AuthService();

export default function UserReset() {
  const navigate = useNavigate();
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<Error | undefined>(undefined);
  const [failedMessage, setFailedMessage] = React.useState<string>("");
  const handleReset = async (input: ResetInput) => {
    try {
        setFailedMessage("");
      setError(undefined);
      setLoading(true);
      const result = await service.sendPassResetMail(input.email);
      console.log(result);
      if (result.isSuccessful) {
        alert(
          "確認用のメールを送信しました。お手数ですがメールの内容を確認して、リセットを行なってください。"
        );
        navigate("/");
        return;
      }
      if (result.isUserNotExists()) {
        setFailedMessage("指定したメールアドレスが存在しません。入力内容をご確認ください。");
      } else if (result.hasError()) {
        setFailedMessage("エラーが発生しました。" + result.errorMessage);
      }
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };
  return (
    <>
      <PageTitle title="パスワードリセット" />
      <Grid container spacing={2}>
        <Grid item xs={12}>
          {error ? (
            <Alert severity="error">
              エラーが発生しました。お手数ですが時間をおいて再度お試しください。
            </Alert>
          ) : (
            <></>
          )}
          {failedMessage !== "" ? (
            <Alert severity="error">{failedMessage}</Alert>
          ) : (
            <></>
          )}
          <ResetForm input={{ email: "" }} onSubmit={handleReset}></ResetForm>
        </Grid>
        <LoadingSpinner message="Loading..." isLoading={loading} />
      </Grid>
    </>
  );
}
