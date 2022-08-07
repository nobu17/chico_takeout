import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import CalendarTable from "./parts/CalendarTable";

export default function Calendar() {
  return (
    <>
      <AdminPage title="予約カレンダー">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <CalendarTable></CalendarTable>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
