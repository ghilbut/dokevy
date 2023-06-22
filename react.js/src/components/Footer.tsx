import * as React from "react";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Link from "@mui/material/Link";
import Grid from "@mui/material/Grid";
import {Facebook, GitHub, Instagram, LinkedIn, Twitter, BrightnessMedium, DensityMedium} from "@mui/icons-material";
import { Box } from "@mui/material";

const Footer = () => {
  return (
      <Box
          component="footer"
          sx={{
              py: 2,
              px: 2,
              mt: 'auto',
              backgroundColor: (theme) =>
                  theme.palette.mode === "light"
                      ? theme.palette.grey[200]
                      : theme.palette.grey[800],
          }}
      >
          <Container maxWidth="xl">
              <Box mt={1}>
                  <Typography variant="body2" color="text.secondary" align="center">
                      <Link href="https://www.linkedin.com/in/ghilbut/" color="inherit">
                          <LinkedIn />
                      </Link>
                      <Link href="https://github.com/ghilbut" color="inherit">
                          <GitHub />
                      </Link>
                  </Typography>
              </Box>
              <Box mt={1}>
                  <Typography variant="body2" color="text.secondary" align="center">
                      {"Copyright © "}
                      <Link color="inherit" href="#">
                          ghilbut
                      </Link>{" "}
                      {new Date().getFullYear()}
                      {"."}
                  </Typography>
              </Box>
          </Container>
      </Box>
  )
}

export default Footer;
