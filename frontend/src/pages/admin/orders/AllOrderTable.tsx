import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import AdminPage from "../../../components/parts/AdminPage";
import OrderTableList from "./parts/OrderTableList";

export default function AllOrderTable() {
  return (
    <>
      <AdminPage title="予約一覧(注文順)">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Typography align="center">最新1,000件のみ表示</Typography>
            <OrderTableList></OrderTableList>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
