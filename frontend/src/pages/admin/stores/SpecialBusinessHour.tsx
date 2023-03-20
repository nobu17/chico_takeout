import { Typography, Grid } from "@mui/material";
import AdminPage from "../../../components/parts/AdminPage";
import SpecialBusinessHourTable from "./parts/SpecialBusinessHourTable";

export default function SpecialBusinessHour() {
  return (
    <>
      <AdminPage title="特別営業時間 編集">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <Typography sx={{ flex: 1 }} variant="caption">
              ※設定した種別に関係なく、同一日の通常営業時間は全て無効になります。
            </Typography>
            <SpecialBusinessHourTable></SpecialBusinessHourTable>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
