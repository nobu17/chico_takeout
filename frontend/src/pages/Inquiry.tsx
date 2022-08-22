import { Typography } from "@mui/material";
import Grid from "@mui/material/Grid";
import PageTitle from "../components/parts/PageTitle";

export default function Inquiry() {
  return (
    <>
      <PageTitle title="お問い合わせ" />
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Typography variant="body1" align="center" paragraph>
            ※お急ぎの用事は、電話かLINE(公式HP参照)でご連絡ください。
            <br />
            TEL:080-6712-2988
            <br />
          </Typography>
          <iframe
            title="問い合わせフォーム"
            src="https://docs.google.com/forms/d/e/1FAIpQLSdIKxlkZnf3qnMvshBMLxQlBqunbPpNtQmupZvUfRn_aL2s7A/viewform?embedded=true"
            width="100%"
            height="898"
            frameBorder={0}
            marginHeight={0}
            marginWidth={0}
          >
            読み込んでいます…
          </iframe>
        </Grid>
      </Grid>
    </>
  );
}
