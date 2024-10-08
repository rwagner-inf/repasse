for env, destinationList := range destinationsMap {
	var computedDestinationList []argocdtypes.ApplicationDestination
	computedDestinationList = append(computedDestinationList, destinationList...)

	appProjectName := fmt.Sprintf("%s-%s-%s", gitlabNamespace.Name, env, gitlabProject.Path)
	resourceVersion := "0"
	existingAppProject, err := pc.ArgoprojV1alpha1().AppProjects(arm.ArgoCDControllerNamespace).Get(ctx, appProjectName, metav1.GetOptions{})
	if err == nil {
		resourceVersion = existingAppProject.ResourceVersion
	} else if !errors.IsNotFound(err) {
		arm.Logger.Error(err, "error finding existing AppProject", "name", appProjectName)
		continue
	}

	var appProjectSpec argocdtypes.AppProjectSpec
	var overrideProjectOk bool
	var overrideProjectSpec argocdtypes.AppProjectSpec
	sourceRepos := []string{gitlabProject.HTTPUrlToRepo}
	sourceRepos = append(sourceRepos, c.UNIVERSAL_REPOS...)
	sort.Strings(sourceRepos)

	if overrideAppProjects != nil {
		overrideProjectSpec, overrideProjectOk = overrideAppProjects[appProjectName]
	}
	if overrideProjectOk {
		arm.Logger.Info("overriding autogenerated AppProjectSpec from ConfigMap", "appProject", appProjectName)
		appProjectSpec = overrideProjectSpec
	} else {
		whitelistedClusterResourcesFinal := []metav1.GroupKind{}
		blacklistedClusterResourcesFinal := []metav1.GroupKind{}
		whitelistedNamespacedResourcesFinal := []metav1.GroupKind{}
		blacklistedNamespacedResourcesFinal := []metav1.GroupKind{}
		var syncWindowsFinal = syncWindowsDefault
		switch arm.CoordinatorMode {
		case c.COORDINATOR_MODE_APPLICATION:
			// these are defined in resource_types.go
			whitelistedClusterResourcesFinal = whitelistedClusterResourcesDefault
			blacklistedClusterResourcesFinal = blacklistedClusterResourcesDefault
			whitelistedNamespacedResourcesFinal = whitelistedNamespacedResourcesDefault
			blacklistedNamespacedResourcesFinal = blacklistedNamespacedResourcesDefault
		case c.COORDINATOR_MODE_INFRASTRUCTURE:??
			// no special allowances at present for infra, it can deploy anything
		}