import React, { useCallback, useEffect, useMemo, useState, memo, useRef } from 'react';
import Grid from '@mui/material/Grid';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import ButtonGroup from '@mui/material/ButtonGroup';
import Box from '@mui/material/Box';
import AddIcon from '@mui/icons-material/Add';
import RemoveIcon from '@mui/icons-material/Delete';
import AddToPhotosIcon from '@mui/icons-material/AddToPhotos';
import Tooltip from '@mui/material/Tooltip';
import { useTranslation } from 'react-i18next';
import { FilterDefinitionFieldsModel } from '../../../../../models/general';
import FilterBuilderField from '../FilterBuilderField';
import { generateKey, buildFieldInitialValue, buildFilterBuilderInitialItems } from '../../utils';
import { LineOrGroup, BuilderInitialValueObject, FieldInitialValueObject, FilterValueObject } from '../../types';

/* eslint-disable @typescript-eslint/no-explicit-any */

interface Props {
  filterDefinitionModel: FilterDefinitionFieldsModel;
  initialValue: BuilderInitialValueObject;
  acceptEmptyLines?: boolean;
  onRemove?: () => void | undefined;
  onChange: (fo: null | FilterValueObject) => void;
}

function FilterBuilder({ filterDefinitionModel, onRemove, onChange, initialValue, acceptEmptyLines }: Props) {
  // Setup translate
  const { t } = useTranslation();
  // States
  const resultsRef = useRef<Record<string, null | any>>({});
  const [groupKey, setGroupKey] = useState(initialValue.group);
  const [items, setItems] = useState<LineOrGroup[]>(initialValue.items);

  // Watch initialValue
  useEffect(() => {
    // Set initial value data
    setGroupKey(initialValue.group);
    setItems(initialValue.items);
    // Clean result reference
    const resRefKeys = Object.keys(resultsRef.current);
    resRefKeys.forEach((key) => {
      if (initialValue.items.findIndex((it) => it.key === key) === -1) {
        // Key isn't found => Need to clean it
        delete resultsRef.current[key];
      }
    });
  }, [initialValue]);

  const localSaveManagementHandler = (newGroupKey: string) => {
    // Create new items
    const newItems = Object.values(resultsRef.current);

    // Check if array is empty and if the accept empty lines isn't enabled
    if (newItems.length === 0 && !acceptEmptyLines) {
      // Don't accepting this
      onChange(null);
      return;
    }

    // Check if array is empty and if the accept empty lines is enabled
    if (newItems.length === 0 && acceptEmptyLines) {
      // Accepting this
      onChange({});
      return;
    }

    // Find one item that doesn't have value
    const withoutValueIndex = newItems.findIndex((v) => !v);
    // Check if there is an item without value
    if (withoutValueIndex !== -1) {
      onChange(null);
      return;
    }

    // Optimization if there is only 1 value
    // Return directly the value
    if (newItems.length === 1) {
      onChange(newItems[0] as any);
      return;
    }

    // Otherwise create and save object
    onChange({ [newGroupKey]: newItems });
  };

  // Add group handler
  const addGroupHandler = useCallback(() => {
    setItems((v) => {
      // Generate key
      const key = generateKey();

      // Create new items list
      const newItems = [
        ...v,
        {
          type: 'group',
          key,
          initialValue: buildFilterBuilderInitialItems(undefined),
        },
      ];
      // Save new result value
      resultsRef.current[key] = null;
      // Save management
      localSaveManagementHandler(groupKey);

      return newItems;
    });
  }, []);

  // Add line handler
  const addLineHandler = useCallback(() => {
    setItems((v) => {
      // Generate key
      const key = generateKey();

      // Create new items list
      const newItems = [...v, { type: 'line', key, initialValue: buildFieldInitialValue(undefined)[0] }];
      // Save new result value
      resultsRef.current[key] = null;
      // Save management
      localSaveManagementHandler(groupKey);

      return newItems;
    });
  }, []);

  // Local remove handler
  const localRemoveHandler = useCallback(
    (key: string) => () => {
      setItems((v) => {
        const newItems = v.filter((it) => it.key !== key);
        // Delete key in result ref
        delete resultsRef.current[key];
        // Save management
        localSaveManagementHandler(groupKey);

        return newItems;
      });
    },
    [],
  );

  // Create all memorized save handlers
  const saveHandlers = useMemo(
    () =>
      items.map((it) => (v: any) => {
        // Save value
        resultsRef.current[it.key] = v;

        // Save management
        localSaveManagementHandler(groupKey);
      }),
    [items],
  );

  return (
    <Box sx={{ display: 'flex' }}>
      <Box
        sx={{
          borderLeftColor: 'text.secondary',
          borderLeftStyle: 'solid',
          borderLeftWidth: '1px',
          margin: '16px 10px 50px 0',
        }}
      />
      <Box sx={{ display: 'block', width: '100%' }}>
        <Box>
          <ButtonGroup size="small" color={items.length === 0 && !acceptEmptyLines ? 'error' : 'primary'}>
            <Button
              onClick={() => {
                setGroupKey('AND');
                // Save management
                localSaveManagementHandler('AND');
              }}
              variant={groupKey === 'AND' ? 'contained' : undefined}
            >
              {t('common.operations.and')}
            </Button>
            <Button
              onClick={() => {
                setGroupKey('OR');
                // Save management
                localSaveManagementHandler('OR');
              }}
              variant={groupKey === 'OR' ? 'contained' : undefined}
            >
              {t('common.operations.or')}
            </Button>
          </ButtonGroup>
          <Tooltip title={<>{t('common.filter.addNewField')}</>}>
            <IconButton onClick={addLineHandler} sx={{ margin: '0 5px' }}>
              <AddIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title={<>{t('common.filter.addNewGroupField')}</>}>
            <IconButton onClick={addGroupHandler} sx={{ margin: '0 5px' }}>
              <AddToPhotosIcon />
            </IconButton>
          </Tooltip>
          {onRemove && (
            <Tooltip title={<>{t('common.filter.deleteGroupField')}</>}>
              <IconButton onClick={onRemove}>
                <RemoveIcon />
              </IconButton>
            </Tooltip>
          )}
        </Box>

        {items.map((it, index) => {
          if (it.type === 'line') {
            return (
              <Box key={it.key} sx={{ display: 'flex', margin: '10px 0 5px 0' }}>
                <Box sx={{ display: 'flex', alignItems: 'center', margin: '-20px 5px 0 0' }}>
                  <Tooltip title={<>{t('common.filter.deleteField')}</>}>
                    <IconButton onClick={localRemoveHandler(it.key)}>
                      <RemoveIcon />
                    </IconButton>
                  </Tooltip>
                </Box>
                <Grid
                  container
                  spacing={1}
                  sx={{
                    // Force line height to always include space for errors
                    minHeight: '72px',
                  }}
                >
                  <FilterBuilderField
                    filterDefinitionModel={filterDefinitionModel}
                    initialValue={it.initialValue as FieldInitialValueObject}
                    onChange={saveHandlers[index]}
                  />
                </Grid>
              </Box>
            );
          }

          return (
            <FilterBuilder
              key={it.key}
              filterDefinitionModel={filterDefinitionModel}
              initialValue={it.initialValue as BuilderInitialValueObject}
              onRemove={localRemoveHandler(it.key)}
              onChange={saveHandlers[index]}
            />
          );
        })}
      </Box>
    </Box>
  );
}

FilterBuilder.defaultProps = {
  onRemove: undefined,
  acceptEmptyLines: false,
};

export default memo(FilterBuilder);