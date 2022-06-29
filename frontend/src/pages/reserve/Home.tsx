import Grid from "@mui/material/Grid";
import PageTitle from "../../components/parts/PageTitle";
import ReserveForm from "./parts/ReserveForm";

export default function ReserveHome() {
  return (
    <>
      <PageTitle title="予約" />
      <Grid container spacing={2}>
        <ReserveForm/>
      </Grid>
    </>
  );
}
