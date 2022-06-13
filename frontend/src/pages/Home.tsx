import Grid from "@mui/material/Grid";
import PageTitle from "../components/parts/PageTitle";

export default function Home() {
  return (
    <>
      <PageTitle title="Root" />
      <Grid container spacing={2}>
        <Grid item xs={12} md={6}></Grid>
      </Grid>
    </>
  );
}
