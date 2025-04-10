import React from 'react';
import Grid from '@mui/material/Grid';
import GridViewItem from './GridViewItem';
import { TodoModel } from '../../../../models/todos';
import CenterLoading from '../../../../components/CenterLoading';
import NoData from '../../../../components/NoData';
import { ConnectionModel } from '../../../../models/general';

interface Props {
  readonly data: ConnectionModel<TodoModel> | undefined;
  readonly loading: boolean;
}

function GridView({ data, loading }: Props) {
  // Check if loading is enabled or not
  if (loading) {
    return <CenterLoading />;
  }

  // Check if data are present
  if (!data || !data.edges) {
    return <NoData />;
  }

  const items: TodoModel[] = data.edges.map((it) => it.node);

  return (
    <Grid container spacing={2} sx={{ padding: '10px' }}>
      {items.map((item) => (
        <GridViewItem item={item} key={item.id} />
      ))}
    </Grid>
  );
}

export default GridView;
