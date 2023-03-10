import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import MessageEditContainer from "./parts/MessageEditContainer";

export default function Messages() {
  return (
    <>
      <AdminPage title="メッセージ 編集">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <MessageEditContainer></MessageEditContainer>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
