import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AnalyzeComponent } from './analyze/analyze.component';
import { DataComponent } from './data/data.component';
import { MainComponent } from './main/main.component';
import { UploadComponent } from './upload/upload.component';

const routes: Routes = [
  { path: '', component: MainComponent },
  { path: 'upload', component: UploadComponent },
  { path: 'data', component: DataComponent },
  { path: 'analyze', component: AnalyzeComponent },
  { path: '**', redirectTo: '', pathMatch: 'full' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
