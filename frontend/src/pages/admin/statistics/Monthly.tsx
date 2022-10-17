import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import MonthlyChartContainer from "./parts/MonthlyChartContainer";

export default function Monthly() {
  return (
    <>
      <AdminPage title="集計(月別)">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <MonthlyChartContainer />
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
