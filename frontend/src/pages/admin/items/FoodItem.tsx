import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import FoodItemTable from "./parts/FoodItemTable"

export default function FoodItem() {
  return (
    <>
      <AdminPage title="アイテム編集(フード)">
        <Grid container spacing={2}>
          <Grid item xs={12}>
              <FoodItemTable></FoodItemTable>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
