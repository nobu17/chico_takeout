import Grid from "@mui/material/Grid";
import PageTitle from "../../components/parts/PageTitle";
import { auth } from "../../libs/firebase/Firebase";
import { uiConfig } from "../../libs/firebase/FirebaseUI";
import StyledFirebaseAuth from "../../libs/firebase/StyledFirebaseAuth";

export default function UserLogin() {
  return (
    <>
      <PageTitle title="ログイン" />
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <StyledFirebaseAuth uiConfig={uiConfig} firebaseAuth={auth} />
        </Grid>
      </Grid>
    </>
  );
}
