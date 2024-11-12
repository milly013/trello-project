import { Routes } from '@angular/router';
import { CreateProjectComponent } from './create-project-component/create-project-component.component';

import { AddTaskComponent } from './add-task/add-task.component';
import { AddUserComponent } from './add-user/add-user.component';
import { RemoveUserComponent } from './remove-user/remove-user.component';
import { RegistrationComponent } from './registration/registration.component';
import { VerificationComponent } from './verification/verification.component';

export const appRoutes: Routes = [
  { path: 'create-project', component: CreateProjectComponent },
  { path: 'add-task',component: AddTaskComponent},
  { path:'add-user',component:AddUserComponent},
  { path:'remove-user'  ,component:RemoveUserComponent},
  { path:'registration',component:RegistrationComponent},
  { path:'verification',component:VerificationComponent}
  
];
