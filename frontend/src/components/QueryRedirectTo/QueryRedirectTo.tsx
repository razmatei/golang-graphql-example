import React from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, DocumentNode, QueryHookOptions, OperationVariables } from '@apollo/client';
import { useParams, Params, Navigate, NavigateProps } from 'react-router-dom';
import CenterLoading, { CenterLoadingProps } from '~components/CenterLoading';
import ErrorsDisplay, { ErrorsDisplayProps } from '~components/ErrorsDisplay';
import NoData, { NoDataTypographyProps } from '~components/NoData';

export interface Props<T, P extends OperationVariables> {
  // Query document to execute.
  query: DocumentNode;
  // Build navigate to function will return a path.
  // If path result is empty or undefined, then "no data" is displayed.
  // If you don't want any "no data" displayed, manage the redirect path
  buildNavigateTo: (params: T | undefined) => string | undefined | null;
  // Build query variables function will return a query variable object.
  buildQueryVariables: (params: Params<string>) => P;
  // Disable center loading subtitle.
  disableCenterLoadingSubtitle?: boolean;
  // Query hook options.
  queryHookOptions?: Omit<QueryHookOptions<T, P>, 'variables'>;
  // No data Typography props.
  noDataTypographyProps?: NoDataTypographyProps;
  // Center loading props.
  centerLoadingProps?: CenterLoadingProps;
  // Graphql Errors props.
  graphqlErrorsProps?: Omit<ErrorsDisplayProps, 'error|errors'>;
  // Navigate props.
  navigateProps?: Omit<NavigateProps, 'to'>;
}

function QueryRedirectTo<T, P extends OperationVariables>({
  query,
  buildQueryVariables,
  buildNavigateTo,
  disableCenterLoadingSubtitle = false,
  queryHookOptions = {},
  noDataTypographyProps = {},
  centerLoadingProps = {},
  graphqlErrorsProps = {},
  navigateProps = {},
}: Props<T, P>) {
  // Setup translate
  const { t } = useTranslation();
  // Get params
  const params = useParams();

  // Build query variables
  const queryVariables: P = buildQueryVariables(params);

  // Query
  const { data, loading, error } = useQuery<T, P>(query, {
    variables: queryVariables,
    fetchPolicy: 'network-only',
    ...queryHookOptions,
  });

  // Check loading
  if (loading) {
    return (
      <CenterLoading
        subtitle={!disableCenterLoadingSubtitle ? t('common.loadingText') : undefined}
        containerBoxSx={{ margin: '15px 0' }}
        {...centerLoadingProps}
      />
    );
  }

  // Check error
  if (error) {
    return <ErrorsDisplay error={error} {...graphqlErrorsProps} />;
  }

  // Build navigate to
  const to = buildNavigateTo(data);

  // Check if to isn't present
  if (!to) {
    return <NoData {...noDataTypographyProps} />;
  }

  return <Navigate to={to} {...navigateProps} />;
}

export default QueryRedirectTo;
