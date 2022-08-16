import Grid from "@mui/material/Grid";
import PageTitle from "../../components/parts/PageTitle";
import LoadingSpinner from "../../components/parts/LoadingSpinner";
import LoginForm, { LoginInput } from "./parts/LoginForm";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../components/contexts/AuthContext";

export default function AdminLogin() {
  const navigate = useNavigate();
  const { signIn, loading } = useAuth();
  const handleSignIn = async (input: LoginInput) => {
    const result = await signIn(input.email, input.password);
    if (result.isSuccessful) {
      navigate("/admin");
    } else {
      alert("ログインに失敗しました。");
    }
  };
  return (
    <>
      <PageTitle title="ログイン" />
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <LoginForm
            input={{ email: "", password: "" }}
            onSubmit={handleSignIn}
          ></LoginForm>
          <LoadingSpinner message="Loading..." isLoading={loading} />
        </Grid>
      </Grid>
    </>
  );
}
