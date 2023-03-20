import { Grid } from "@mui/material";
import AdminPage from "../../../components/parts/AdminPage";
import BusinessHourTable from "./parts/BusinessHourTable";

export default function BusinessHour() {
  return (
    <>
      <AdminPage title="営業時間 編集">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <BusinessHourTable></BusinessHourTable>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
