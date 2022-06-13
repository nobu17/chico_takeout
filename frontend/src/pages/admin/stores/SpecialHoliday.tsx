import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import SpecialHolidayTable from "./parts/SpecialHolidayTable"

export default function SpecialHoliday() {
  return (
    <>
      <AdminPage title="特別休日 編集">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <SpecialHolidayTable></SpecialHolidayTable>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
