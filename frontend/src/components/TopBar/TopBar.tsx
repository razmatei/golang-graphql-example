import React from 'react';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Avatar from '@mui/material/Avatar';
import { useTranslation } from 'react-i18next';
import IconToggleColorMode from '../theming/IconToggleColorMode';

function TopBar() {
  // Setup translate
  const { t } = useTranslation();

  return (
    <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
      <Toolbar variant="dense">
        <Avatar src="/logo.png" />
        <Typography variant="h6" component="div" sx={{ flexGrow: 1, marginLeft: '10px' }}>
          {t('common.mainTitle')}
        </Typography>
        <IconToggleColorMode />
      </Toolbar>
    </AppBar>
  );
}

export default TopBar;