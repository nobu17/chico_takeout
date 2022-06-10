import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import SpecialBusinessHourTable from "./parts/SpecialBusinessHourTable"

export default function SpecialBusinessHour() {
  return (
    <>
      <AdminPage title="特別営業時間 編集">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <SpecialBusinessHourTable></SpecialBusinessHourTable>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
