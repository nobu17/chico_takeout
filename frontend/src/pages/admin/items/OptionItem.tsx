import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import OptionItemTable from "./parts/OptionItemTable";

export default function OptionItem() {
  return (
    <>
      <AdminPage title="オプションアイテム編集">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <OptionItemTable />
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
