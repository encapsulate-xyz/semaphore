import { test, expect } from './fixtures';

test('run task from demo project', async ({ page, login, project }) => {
  await login(true);
  await project('task_runner', true);

  await page.getByText('Project Test created').waitFor();

  await page.getByTestId('sidebar-templates').click();
  
  await page.getByRole('link', { name: 'Build demo app' }).click();

  await page.getByTestId('template-run').click();

  await page.getByTestId('newTaskDialog').getByRole('textbox', { name: 'Message (Optional)' }).fill('Test');

  await page.getByTestId('newTaskDialog').getByTestId('editDialog-save').click();

  await page.getByTestId('task-rawLog').waitFor({timeout: 100000});

  await expect(page.getByTestId('task-status')).toHaveText('Success');
});

test('stop task on waiting', async ({ page, login, project }) => {
  await login(true);
  await project('task_runner', true);

  await page.getByText('Project Test created').waitFor();

  await page.getByTestId('sidebar-templates').click();
  
  await page.getByRole('link', { name: 'Build demo app' }).click();

  await page.getByTestId('template-run').click();

  await page.getByTestId('newTaskDialog').getByRole('textbox', { name: 'Message (Optional)' }).fill('Test');

  await page.getByTestId('newTaskDialog').getByTestId('editDialog-save').click();

  await page.getByRole('dialog').getByRole('button', { name: 'Stop' }).click();

  await page.getByTestId('task-rawLog').waitFor();

  await expect(page.getByTestId('task-status')).toHaveText('Stopped');

});

test('stop task on cloning', async ({ page, login, project }) => {
    await login(true);
    await project('task_runner', true);
  
    await page.getByText('Project Test created').waitFor();
  
    await page.getByTestId('sidebar-templates').click();
    
    await page.getByRole('link', { name: 'Build demo app' }).click();
  
    await page.getByTestId('template-run').click();
  
    await page.getByTestId('newTaskDialog').getByRole('textbox', { name: 'Message (Optional)' }).fill('Test');
  
    await page.getByTestId('newTaskDialog').getByTestId('editDialog-save').click();

    await page.getByRole('dialog').getByText('Get current commit hash').waitFor();

    await page.getByRole('dialog').getByRole('button', { name: 'Stop' }).click();
  
    await page.getByTestId('task-rawLog').waitFor();

    await expect(page.getByTestId('task-status')).toHaveText('Stopped');
  
  });

  test('stop task on running', async ({ page, login, project }) => {
    await login(true);
    await project('task_runner', true);
  
    await page.getByText('Project Test created').waitFor();
  
    await page.getByTestId('sidebar-templates').click();
    
    await page.getByRole('link', { name: 'Build demo app' }).click();
  
    await page.getByTestId('template-run').click();
  
    await page.getByTestId('newTaskDialog').getByRole('textbox', { name: 'Message (Optional)' }).fill('Test');
  
    await page.getByTestId('newTaskDialog').getByTestId('editDialog-save').click();
  
    await page.getByRole('dialog').getByText('TASK [Gathering Facts] *********************************************************').waitFor();

    await page.getByRole('dialog').getByRole('button', { name: 'Stop' }).click();
  
    await page.getByTestId('task-rawLog').waitFor();

    await expect(page.getByTestId('task-status')).toHaveText('Stopped');
    
  });