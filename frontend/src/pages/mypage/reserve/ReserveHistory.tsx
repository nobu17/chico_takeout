import Grid from "@mui/material/Grid";
import MyPage from "../parts/MyPage";
import ReserveHistoryTable from "./parts/ReserveHistoryTable";

export default function MyReserveHistory() {
  return (
    <>
      <MyPage title="注文履歴">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <ReserveHistoryTable></ReserveHistoryTable>
          </Grid>
        </Grid>
      </MyPage>
    </>
  );
}
